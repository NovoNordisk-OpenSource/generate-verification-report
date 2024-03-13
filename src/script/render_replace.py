import sys

# TODO: Add exception handling
def main(argv):
    # 1. Check for the arg pattern:
    #   python3 render_replace.py -render <filepath> -template <filepath> -placeholder <placeholder identifier>
    #   e.g. 
    #       argv[0] is '-render'
    #       argv[1] is './testResultsHtml.txt'
    #       argv[2] is '-template'
    #       argv[3] is './VerificationReportTemplate.html'
    #       argv[4] is '-placeholder'
    #       argv[5] is '<var>TESTCASE_RESULTS</var>'
    if len(argv) == 6 and argv[0] == '-render' and argv[2] == '-template' and argv[4] == '-placeholder':
        # Read the render HTML from given file
        render_data = ''
        with open(argv[1], 'rt') as f:
            render_data = f.read()
            f.close()

        # Read the template file content
        template_data = ''
        with open(argv[3], 'rt') as f:
            template_data = f.read()
            # Replace placeholder occurrence in the template with render HTML
            template_data = template_data.replace(argv[5], render_data)
            f.close()

        # Write back the template with the render HTML embedded
        with open(argv[3], 'wt') as f:
            f.write(template_data)
            f.close()

if __name__ == "__main__":
   main(sys.argv[1:])