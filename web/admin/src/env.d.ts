declare module '@kangc/v-md-editor' {
  import { App, Plugin } from 'vue'
  const VMdEditor: Plugin & { use: (theme: any) => void }
  export default VMdEditor
}

declare module '@kangc/v-md-editor/lib/theme/github.js' {
  const githubTheme: any
  export default githubTheme
}

declare module '@kangc/v-md-editor/lib/style/base-editor.css'
declare module '@kangc/v-md-editor/lib/theme/style/github.css'
