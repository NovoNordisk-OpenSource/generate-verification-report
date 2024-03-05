# Take inputs: 1) The directory where the feature files are located.

# 1) Parse all the requirements files (feature files) and create a mapping from the Feature:<text> to the unique requirement id found
#    in the top of the feature file.

# 2) Emit a dictionary mapping from the feature description to the unique ID so it can be used as a lookup table. Also do this mapping
#    for the feature description to the Feature file name (so it can be linked to in the verification report).

import sys
import glob

# Template for generating mapping lookup script
TEMPLATE = '''
import sys

# Define mapping dictionary (key: feature name, value: feature tags)
mapping = {mapping}

def main(argv):
    # 1. Check for the arg pattern:
    #   python3 requirements_id_mapping_lookup.py <feature name>
    #   e.g. args[0] is 'Evaluation of Inference Unit with multiple sub units'
    if len(argv) == 1 and argv[0] != '':
        try:
            print(mapping[argv[0]])   
        except (KeyError):
            print('')
    else:
        print('MAPPING_ERROR')


if __name__ == "__main__":
    main(sys.argv[1:])

'''

# List of tags that are being removed when rendering the list of requirements, leaving only the unique ID(s)
RESERVED_TAGS = ["URS", "GxP", "non-GxP", "CA"]

class Feature:
    def __init__(self, name, tags):
        self.name = name
        self.tags = tags

def extract_features(lines):
    features = []
    for i, line in enumerate(lines):
        if '@URS' in line:
            fo = Feature(None, tags=clean_tags(line))
            for j in range(i + 1, len(lines)):
                feature_line = lines[j]
                if 'Feature:' in feature_line:
                    feature_name = feature_line.split('Feature:')[1].strip()
                    fo.name = feature_name
                    features.append(fo)
                    break

    return features

def remove_values_from_string(string, values):
    for value in values:
        string = string.replace(value, '')
    return string

def clean_tags(tags) -> list:
    tags = remove_values_from_string(tags, RESERVED_TAGS)
    tags = tags.replace('@', '')
    tags = tags.strip()
    return tags

def main(argv):
    # 1. Check for the arg pattern:
    #   python3 extract_requirements_name_to_id_mapping.py -folder <folderpath>
    #   e.g. args[0] is '-folder' and args[1] is './requirements'
    if len(argv) == 2 and argv[0] == '-folder':
        # Find all feature files in the folder and subfolders
        path = r'%s/**/*.feature' % argv[1]
        files = glob.glob(path, recursive=True)

        # Creating mapping dictionary
        mapping = {}

        # Create a mapping for each unqiue feature id to the feature description
        for i, file in enumerate(files):
            #print(file)
            with open(file, mode='r', encoding='utf-8') as file_reader:
                lines = file_reader.read().split('\n')
                features = extract_features(lines)
                file_reader.close()
                for feature in features:
                    mapping[feature.name] = feature.tags

        # Print the mapping
        #print(mapping)

        # # Write the mapping script to a file
        # with open('requirements_id_mapping_lookup.py', mode='w', encoding='utf-8') as file_writer:
        #     file_writer.write(TEMPLATE.replace('{mapping}', str(mapping)))
        #     file_writer.close()

        # Emit the mappings as stdout
        print(str(mapping))


if __name__ == "__main__":
    main(sys.argv[1:])

