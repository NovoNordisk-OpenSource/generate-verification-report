import sys
import glob
import os
import subprocess

# List of tags that are being removed when rendering the list of requirements, leaving only the unique ID(s)
RESERVED_TAGS = ["URS", "GxP", "non-GxP", "CA"]
class Feature:
    def __init__(self, name, description, tags):
        self.name = name
        self.description = description
        self.tags = tags

def extract_features(lines):
    features = []
    for i, line in enumerate(lines):
        if '@URS' in line:
            fo = Feature(None, description=[], tags=line)
            for j in range(i + 1, len(lines)):
                feature_line = lines[j]
                if 'Feature:' in feature_line:
                    feature_name = feature_line.split('Feature:')[1].strip()
                    feature_description = lines[j + 1:j + 5]
                    fo.name = feature_name
                    fo.description = [desc.strip() for desc in feature_description if desc.strip()]
                    features.append(fo)
                    break

    return features

def extract_last_modified_commit_hash(filepath, branch):
    # git log <remote branch> -n 1 --pretty=format:%H -- <filepath>
    result = subprocess.run(['git', 'log', branch, '-n', '1', '--pretty=format:%H', '--', filepath], stdout=subprocess.PIPE)
    return result.stdout.decode()

def extract_last_modified_commit_hash_timestamp(commit_hash):
    # git show -s --format=%cd --date=format:'%Y-%m-%d %H:%M:%S' <commit_hash>
    result = subprocess.run(['git', 'show', '-s', '--format=%cd', "--date=format:'%Y-%m-%d %H:%M:%S'", commit_hash], stdout=subprocess.PIPE)
    return result.stdout.decode()

def remove_values_from_string(string, values):
    for value in values:
        string = string.replace(value, '')
    return string

def clean_tags(tags) -> list:
    tags = remove_values_from_string(tags, RESERVED_TAGS)
    tags = tags.replace('@', '')
    tags = tags.strip()
    return tags

def render_tags(tags) -> str:
    tags = clean_tags(tags)
    tags = tags.split(' ')
    tags = [f'<kbd>{tag}</kbd>' for tag in tags]
    return ''.join(tags)

# TODO: Add exception handling
def main(argv):
    # Guard clause, no arguments provided
    if len(argv) == 0:
        print("No arguments provided")
        exit(1)
    # Guard clause, too few arguments provided
    if len(argv) < 10:
        print("Not all required arguments provided")
        exit(1)

    # 1. Check for the arg pattern:
    #   python3 get_requirements.py -folder <filepath> -branch <remote branch> -organization <organization> -project <project> -repository <repository>
    #   e.g. 
    #       argv[0] is '-folder'
    #       argv[1] is './../features'
    #       argv[2] is '-branch'
    #       argv[3] is 'origin/release/service1'
    #       argv[4] is '-organization'
    #       argv[5] is 'novonordiskit'
    #       argv[6] is '-project'
    #       argv[7] is 'Data Management and Analytics'
    #       argv[8] is '-repository'
    #       argv[9] is 'QMS-TEMPLATE'
    if len(argv) == 10 and argv[0] == '-folder' and argv[2] == '-branch' and argv[4] == '-organization' and argv[6] == '-project' and argv[8] == '-repository':
	
        # Render all feature descriptions
        # Find all .feature files in the folder and subfolders
        path = r'%s/**/*.feature' % argv[1]
        files = glob.glob(path, recursive=True)

        # Get the current branch
        branch = argv[3]
        organization = argv[5]
        project = argv[7]
        repository = argv[9]

        # URL encode the project name
        project = project.replace(" ", "%20")
        
        # Render the table header and the table body element
        print('''<figure>
    <table>
        <thead>
            <tr>
                <th scope="col">#</th>
                <th scope="col">Requirement ID</th>
                <th scope="col">Descriptive title</th>
                <th scope="col">Version</th>
                <th scope="col">Last modified</th>
                <th scope="col">File</th>
            </tr>
        </thead>
        <tbody>''')

        count_features = 0
        for file in files:
            last_modified_commit_hash = extract_last_modified_commit_hash(file, branch)
            last_modified_commit_hash_timestamp = extract_last_modified_commit_hash_timestamp(last_modified_commit_hash).strip().replace("'", "")

            with open(file, mode='r', encoding='utf-8') as file_reader:
                lines = file_reader.read().split('\n')
                features = extract_features(lines)
                # Extract the path to the feature file, e.g.:
                # /requirements/features/urs/functionality1.feature
                repository_file_path = os.path.abspath(file).replace(os.getcwd(), "")
                # Create link to path for file, e.g.:
                # https://dev.azure.com/novonordiskit/Data%20Management%20and%20Analytics/_git/QMS-TEMPLATE/commit/d78d1bf6bd41b07f654c6b8178fb85b4490853f3?path=/requirements/features/urs/reverse-string-feat.feature
                repository_file_link = f'https://dev.azure.com/{organization}/{project}/_git/{repository}/commit/{last_modified_commit_hash}?path={repository_file_path}'

                for feature in features:
                    count_features += 1
                    print(f'''            <tr>
                <th scope="row">{count_features}</th>
                <td>{render_tags(feature.tags)}</td>
                <td>{feature.name}</td>
                <td>{last_modified_commit_hash}</td>
                <td>{last_modified_commit_hash_timestamp}</td>
                <td><a href="{repository_file_link}" target="_blank">{repository_file_path}</a></td>
            </tr>''')

        # Render the table body close elements
        print('''        </tbody>
    </table>
</figure>''')

if __name__ == "__main__":
   main(sys.argv[1:])