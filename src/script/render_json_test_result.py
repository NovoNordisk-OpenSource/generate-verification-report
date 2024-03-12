import sys
import json 
import glob
from datetime import datetime

class Testresult:
    def __init__(self, name, status, start, stop, uuid, historyId, fullName, labels, statusDetails =None, steps =None, **kwargs):
        self.name = name
        self.status = status
        self.statusdetails = statusDetails
        self.steps = steps
        self.start = start
        self.stop = stop
        self.uuid = uuid
        self.historyid = historyId
        self.fullname = fullName 
        self.labels = labels
        self.count = 0

    @classmethod
    def from_object(cls, json_object):
        return cls(**json_object)

    # TODO: Add exception handling
    def __repr__(self):
        # Prepare values for rendering
        time_elapsed = (
            datetime.fromtimestamp(self.stop/1000.0) - 
            datetime.fromtimestamp(self.start/1000.0)
            ).total_seconds() # Time elapsed in seconds for test to execute
        time_executed = datetime.fromtimestamp(self.start/1000.0).strftime('%Y-%m-%d %H:%M:%S') # Time test was executed
        # Extract tags and features from labels list (TODO: This needs to be refactored into input in a future interface)
        #tags = [x['value'] for x in self.labels if x['name'] == 'tag']
        features = [x['value'] for x in self.labels if x['name'] == 'feature']
        na_tag = 'N/A'

        # Lookup all features' unique tag from feature name in the mapping dictionary (depends on provided output from 'extract_requirements_name_to_id_mapping.py')
        features_tags = []
        try:
            for feature in features:
                # Append the feature tag to the list only if not empty
                feature_tag = mapping[feature]
                if feature_tag != '':
                    features_tags.append(feature_tag)
        except (KeyError):
            pass
        
        # Create rendering
        return f'''            <tr>
                <th scope="row">{self.count}</th>
                <td>{self.name}</td>
                <td>Automatic</td>
                <td>Pipeline</td>
                <td>{time_executed}</td>
                <td class="{self.status}">{self.status}</td>
                <td>{'<kbd>' + '</kbd><kbd>'.join(features_tags) + '</kbd>' if features_tags else ''.join(na_tag)}</td>
            </tr>'''

def render_header():
    return '''<figure>
    <table>
        <thead>
            <tr>
                <th scope="col">#</th>
                <th scope="col">Test case title</th>
                <th scope="col">Method</th>
                <th scope="col">Executed by</th>
                <th scope="col">Executed at</th>
                <th scope="col">Conclusion</th>
                <th scope="col">Trace to requirements</th>
            </tr>
        </thead>
        <tbody>'''

def render_footer():
    return '''        </tbody>
    </table>
</figure>'''

# TODO: Add exception handling
def main(argv):
    # METHOD 1: DEPRECATED (2023-04-24) - TO BE REMOVED
    # 1. Check for the arg pattern:
    #   python3 render_json_test_result.py -file <filepath>
    #   e.g. args[0] is '-file' and args[1] is './results/cf0355e3-be5f-4d57-b103-fc751059b394-result.json'
    if len(argv) == 2 and argv[0] == '-file':
        # Create rendering for the test result
        with open(argv[1], 'r') as f:
            json_data = json.load(f)
            testresult = Testresult.from_object(json_data)
            f.close()
            print(testresult)

    # 2. Check for the arg pattern:
    #   python3 render_json_test_result.py -folder <folderpath> -mapping <mappingfilepath>
    #   e.g. args[0] is '-folder' and args[1] is './results' and args[2] is '-mapping' and args[3] is './mapping.dict'
    if len(argv) == 4 and argv[0] == '-folder' and argv[2] == '-mapping':
        # Set the mapping dictionary with the contents of the provided mapping file
        with open(argv[3], mode='r', encoding='utf-8') as file_reader:
            # Define a globally accessible mapping dictionary (key: feature name, value: feature tags)
            global mapping
            mapping = eval(file_reader.read())
            file_reader.close()

        # Find all test result json files in the folder and subfolders
        path = r'%s/**/*-result.json' % argv[1]
        files = glob.glob(path, recursive=True)

        # Render the table header and body open elements
        print(render_header())

        # Create rendering for each test result
        for i, file in enumerate(files):
            with open(file, 'r') as f:
                json_data = json.load(f)
                testresult = Testresult.from_object(json_data)
                testresult.count = i + 1
                f.close()
                print(testresult)

        # Render the table body close elements
        print(render_footer())

if __name__ == "__main__":
    main(sys.argv[1:])