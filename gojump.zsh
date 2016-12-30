# change pwd hook
gojump_chpwd() {
    gojump --add "$(pwd)" >/dev/null &!
}

typeset -gaU chpwd_functions
chpwd_functions+=gojump_chpwd

# default gojump command
j() {
	if [[ ${1} == "" ]]; then
		cd
		return
	fi
    if [[ ${1} == -* ]] && [[ ${1} != "--" ]]; then
        gojump ${@}
        return
    fi

    setopt localoptions noautonamedirs
    local output="$(gojump ${@})"
    if [[ -d "${output}" ]]; then
        if [ -t 1 ]; then  # if stdout is a terminal, use colors
                echo -e "\\033[31m${output}\\033[0m"
        else
                echo -e "${output}"
        fi
        cd "${output}"
    else
        echo "gojump: directory '${@}' not found"
        echo "${output}"
        echo "Try \`gojump --help\` for more information."
        false
    fi
}
