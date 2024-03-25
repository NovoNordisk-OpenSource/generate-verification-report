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
    if len(argv) < 8:
        print("Not all required arguments provided")
        exit(1)

    # 1. Check for the arg pattern:
    #   python3 render_requirements_for_github.py -folder <filepath> -branch <remote branch> -organization <organization> -repository <repository>
    #   e.g. 
    #       argv[0] is '-folder'
    #       argv[1] is './../features'
    #       argv[2] is '-branch'
    #       argv[3] is 'origin/release/service1'
    #       argv[4] is '-organization'
    #       argv[5] is 'nn-dma'
    #       argv[6] is '-repository'
    #       argv[7] is 'generate-verification-report'
    if len(argv) == 8 and argv[0] == '-folder' and argv[2] == '-branch' and argv[4] == '-organization' and argv[6] == '-repository':
	
        # Render all feature descriptions
        # Find all .feature files in the folder and subfolders
        path = r'%s/**/*.feature' % argv[1]
        files = glob.glob(path, recursive=True)

        # Get the current branch
        branch = argv[3]
        organization = argv[5]
        repository = argv[7]
        
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
                # https://github.com/nn-dma/generate-verification-report/blob/5ef02fe1e00c1dbf9d924ce9717af85e9d83ae44/test/integration/requirements/urs/uppercase-string.feature
                repository_file_link = f'https://github.com/{organization}/{repository}/blob/{last_modified_commit_hash}/{repository_file_path}'

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