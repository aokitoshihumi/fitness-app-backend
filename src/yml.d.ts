// wrangler.jsonc の rules で .yml を文字列としてimportできるようにしている
declare module "*.yml" {
  const content: string;
  export default content;
}
