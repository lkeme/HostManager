# README

## Doc
- Wails: https://wails.io/zh-Hans/docs/gettingstarted/firstproject/
- Element Plus: https://element-plus.org/zh-CN/guide/installation.html
- Vue Router: https://router.vuejs.org/zh/
- Vuex: https://next.vuex.vuejs.org/zh/
- Axios: https://axios-http.com/docs/intro
- Vite: https://cn.vitejs.dev/
- TypeScript: https://www.typescriptlang.org/zh/
- Go: https://golang.org/doc/
- GoDoc: https://pkg.go.dev/
- GoFrame: https://goframe.org/
- GoFrame Doc: https://goframe.org/pages/viewpage.action?pageId=1114397
- Tailwind CSS: https://tailwindcss.com/docs
- Naive UI: https://www.naiveui.com/zh-CN/light/components/tree

## xicons

> https://github.com/07akioni/xicons#installation

> https://xicons.org/#/

```
# Icons Installation For vue3
npm i -D @vicons/fluent
npm i -D @vicons/ionicons4
npm i -D @vicons/ionicons5
npm i -D @vicons/antd
npm i -D @vicons/material
npm i -D @vicons/fa # font awesome
npm i -D @vicons/tabler
npm i -D @vicons/carbon
# Icon Utils Installation For vue3
npm i -D @vicons/utils
```

## About

This is the official Wails Vue-TS template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.
