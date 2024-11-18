set_skip_usr_input=0
set_reboot=0

get_usr_input()
{
	echo $1
	if [ $set_skip_usr_input -eq 0 ]; then
		read -r userInput
	else
		userInput=Y
	fi

	if [ "$userInput" != 'Y' ] && [ "$userInput" != 'y' ]; then
		return 0
	else
		return 1
	fi
}

usage()
{
	echo "Usage: $(basename $0) [OPTION]"
	echo "    -s, --skip   设置跳过询问"
	echo "    -r, --run    正常执行"
	echo "    -h, --help   显示帮助信息"
}

while true
do
	case "$1" in
		-s|--skip)
			set_skip_usr_input=1
{{range .ShellFuns}}			{{.}}
{{ end }}
			exit 0
			;;
		-r|--run)
{{range .ShellFuns}}			{{.}}
{{ end }}
			exit 0
			;;
		-h|--help)
			usage
			exit $?
			;;
		*)
			echo -e "Need Correct Arguments!\n"
			usage
			exit $LA_ERR
			;;
	esac
done
