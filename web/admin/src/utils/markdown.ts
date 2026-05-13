import MarkdownIt from 'markdown-it'
import TurndownService from 'turndown'

const md = new MarkdownIt({
  html: true,
  breaks: true,
  linkify: true,
})

const turndown = new TurndownService({
  headingStyle: 'atx',
  codeBlockStyle: 'fenced',
  emDelimiter: '*',
})

export function markdownToHtml(markdown: string): string {
  if (!markdown) return ''
  return md.render(markdown)
}

export function htmlToMarkdown(html: string): string {
  if (!html) return ''
  return turndown.turndown(html)
}
