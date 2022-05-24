export type Service = {
  name: string
  method: 'tcp'|'udp'|'ping'
  status: 'offline'|'online'|'pending'|'error'
  delay: number
  address: string
  port: number
}