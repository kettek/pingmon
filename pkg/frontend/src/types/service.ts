export type Service = {
  name: string
  status: 'offline'|'online'|'pending'
  address: string
}