use t;

sub run($code) { post-solution(:$code)<Out> }

is run('say %*ENV'), '{}', 'Environment';

is run('say $*KERNEL.hostname'), 'code-golf', 'Hostname';

done-testing;
