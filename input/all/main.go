package all

import (
	_ "dev.sum7.eu/genofire/logmania/input/journald_json"
	_ "dev.sum7.eu/genofire/logmania/input/logrus"
	_ "dev.sum7.eu/genofire/logmania/input/syslog"
	_ "dev.sum7.eu/genofire/logmania/input/webhook"
	_ "dev.sum7.eu/genofire/logmania/input/webhook/circleci"
	_ "dev.sum7.eu/genofire/logmania/input/webhook/git"
	_ "dev.sum7.eu/genofire/logmania/input/webhook/grafana"
)
