import sys

def main(argv):
    if len(argv) == 4 and argv[0] == '-orgrepo' and argv[2] == '-issue':
        org_repo = argv[1]
        issue_no = argv[3]

        it_change_issue_link = f'https://github.com/{org_repo}/issues/{issue_no}'

        if len(argv) == 0:
            print(f"<kbd>!MISSING!</kbd>")
        else:
            print(
                f"<kbd><a href=\"{it_change_issue_link}\">#{issue_no}</a></kbd>"
            )

if __name__ == "__main__":
   main(sys.argv[1:])
