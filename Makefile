in = ""
now = $(shell date +%Y%m%d)
.PHONY: reset commit push
reset:
	@git fetch origin
	@git reset --hard origin/main
	@git clean -f -d
	@echo "=========================="
	@echo "Local code reset complete"
	@echo "=========================="
commit:
	@git config --global user.email "631821745@qq.com"
	@git config --global user.name "Gnoloayoul"
	@git add .
	@git commit -s -m "$(in) -$(now)"
	@echo "======================================="
	@echo "Git push ready, please exec [git push]"
	@echo "======================================="
push: commit
	@git push origin main
	@echo "============================"
	@echo "JGEBCamp git push accomplish"
	@echo "============================"