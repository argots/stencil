# stencil
A package manager for static files and templates

# todo

- [ ] `stencil pull file.tpl` should run the file as a go template.
- [ ] Add template function `file(localPath, githubURL)` which copies file at githubURL.
- [ ] Add 3-way merge if git pull brings newer file and local file also modified.
- [ ] Add template function `stencil(localPath, githubURL)` which also treats file as a template.
- [ ] `stencil pull github.com/...` souuld fetch from public github (standard github url)
- [ ] Update `file(..)` to support relative github URLs
- [ ] `stencil pull github.com/...` souuld fetch from private github using ssh settings
- [ ] Add template variables `bool("name", "prompt")` and `var("name")`
- [ ] Add template variables `string("name", "prompt")` and `var("name")`
- [ ] `stencil pull` should pull latest versions of all URLs that were already pulled.
- [ ] Add nested templates support: `import(otherFile)`
- [ ] Add nested pull support `pull(args)`
- [ ] Add ability to modify variables `stencil -var XYZ value`
- [ ] Add ability to use keyrings for secrets `secret("name", "prompt")` and `var("name")`
- [ ] Add ability to work with file patches inserted using markers
- [ ] Deal with diamond dependencies?
- [ ] Unsafe shell exec?