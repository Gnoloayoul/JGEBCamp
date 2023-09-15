in = ""
now = $(shell date +%Y-%m-%d)
.PHONY: reset gitCommit test
reset:
	@git fetch origin
	@git reset --hard origin/main
	@git clean -f -d
	@echo "=========================="
	@echo "Local code reset complete"
	@echo "=========================="
gitCommit:
	@git config --global user.email "631821745@qq.com"
	@git config --global user.name "Gnoloayoul"
	@git add .
	@git commit -s -m "$(in)-$(now)"
	@echo "=========================="
	@echo "Git push ready"
	@echo "=========================="