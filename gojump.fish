# enable tab completion
complete -x -c j -a '(gojump -compl (commandline -t))'

# change pwd hook
function __chpwd_gojump --on-variable PWD --description 'handler of changing $PWD'
    if not status --is-command-substitution ; and status --is-interactive
        status --is-command-substitution; and return
        gojump -add $PWD &
    end
end

# default gojump command
function j
    set -l output (gojump $argv)
    if test -d $output
        set_color red
        echo $output
        set_color normal
        cd $output
    else
        echo "gojump: directory '"$argv"' not found"
    end
end

# open gojump results in file browser
function jo
    set -l output (gojump $argv)
    if test -d $output
        xdg-open $output
    else
        echo "gojump: directory '"$argv"' not found"
    end
end
