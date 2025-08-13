# D:\dev\go\foo_lang_v2\syntax\vscode
# syntax\vscode\package.json change version
# cd syntax/vscode && npm run compile
# cd syntax/vscode && npx vsce package
#cd syntax/vscode && code --install-extension foo-lang-1.1.3.vsix --force
uninstall:
	code --uninstall-extension foolang.foo-lang

install:
# 	cd syntax\vscode\package.json change version
# 	cd lsp && go build -o foo-lsp.exe .
	cd lsp && go install .
	cd syntax/vscode && bun run compile \
	&& bunx vsce package \
	&& code --install-extension foo-lang-1.1.5.vsix --force
