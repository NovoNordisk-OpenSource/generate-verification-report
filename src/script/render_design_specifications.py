import sys
import glob
import os
import subprocess


def extract_design_specification_tags(lines):
    tags = []
    readTags = False
    for i, line in enumerate(lines):
        if 'tags:' in line:
            # tags:
            readTags = True
        if readTags:
            #   - service_now_integration
            if '- ' in line:
                tags.append(line.split('- ')[1].strip())
            if '---' in line: # This marks the end of a tag section, so do not process the rest of the file.
                # ---
                break
    return tags

def extract_last_modified_commit_hash(filepath, branch):
    # git log <remote branch> -n 1 --pretty=format:%H -- <filepath>
    result = subprocess.run(['git', 'log', branch, '-n', '1', '--pretty=format:%H', '--', filepath], stdout=subprocess.PIPE)
    return result.stdout.decode()

def extract_last_modified_commit_hash_timestamp(commit_hash):
    # git show -s --format=%cd --date=format:'%Y-%m-%d %H:%M:%S' <commit_hash>
    result = subprocess.run(['git', 'show', '-s', '--format=%cd', "--date=format:'%Y-%m-%d %H:%M:%S'", commit_hash], stdout=subprocess.PIPE)
    return result.stdout.decode()

# TODO: Add exception handling
def main(argv):
    # Guard clause, no arguments provided
    if len(argv) == 0:
        print("No arguments provided")
        exit(1)
    # Guard clause, too few arguments provided
    if len(argv) < 6:
        print("Not all required arguments provided")
        exit(1)

    # 2. Check for the arg pattern:
    #   python3 render_requirements_for_github.py -folder <filepath> -branch <remote branch> -repository <repository>
    #   e.g. 
    #       argv[0] is '-folder'
    #       argv[1] is './../features'
    #       argv[2] is '-branch'
    #       argv[3] is 'origin/release/v1'
    #       argv[4] is '-repository'
    #       argv[5] is 'nn-dma/generate-verification-report'
    if len(argv) == 6 and argv[0] == '-folder' and argv[2] == '-branch' and argv[4] == '-repository':
        # Render all design specifications
        # Find all .md files in the folder and subfolders
        path = r'%s/**/*.md' % argv[1]
        files = glob.glob(path, recursive=True)

        # Get the current branch
        branch = argv[3]
        repository = argv[5]

        # Render the table header and the table body element
        print('''<figure>
    <table>
        <thead>
            <tr>
                <th scope="col">#</th>
                <th scope="col">Design specification</th>
                <th scope="col">Version</th>
                <th scope="col">Last modified</th>
                <th scope="col">Trace to requirements</th>
            </tr>
        </thead>
        <tbody>''')

        count_design_specifications = 0
        for file in files:
            #print(f"File: {file}")
            last_modified_commit_hash = extract_last_modified_commit_hash(file, branch)
            last_modified_commit_hash_timestamp = extract_last_modified_commit_hash_timestamp(last_modified_commit_hash).strip().replace("'", "")

            with open(file, mode='r', encoding='utf-8') as file_reader:
                lines = file_reader.read().split('\n')
                design_specification_tags = extract_design_specification_tags(lines)
                # Extract the path to the markdown file, e.g.:
                # /system_documentation/docs/design/ReverseString/index.md
                repository_file_path = os.path.abspath(file).replace(os.getcwd(), "")
                # Create link to path for file, e.g.:
                # https://github.com/nn-dma/generate-verification-report/blob/5ef02fe1e00c1dbf9d924ce9717af85e9d83ae44/test/integration/requirements/urs/uppercase-string.feature
                repository_file_link = f'https://github.com/{repository}/blob/{last_modified_commit_hash}{repository_file_path}'

                count_design_specifications += 1
                print(f'''            <tr>
                <th scope="row">{count_design_specifications}</th>
                <td><a href="{repository_file_link}" target="_blank">{repository_file_path}</a></td>
                <td>{last_modified_commit_hash}</td>
                <td>{last_modified_commit_hash_timestamp}</td>
                <td>{ '<kbd>' + '</kbd><kbd>'.join(design_specification_tags) + '</kbd>' if design_specification_tags else 'N/A' }</td>
            </tr>''')

        # Render the table body close elements
        print('''        </tbody>
    </table>
</figure>''')

if __name__ == "__main__":
   main(sys.argv[1:])