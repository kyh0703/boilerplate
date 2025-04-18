import { registerAs } from '@nestjs/config'

export default registerAs('app', () => ({
  port: process.env.HTTP_PORT || 8000,
  cookieKey: process.env.COOKIE_KEY,
  frontURL: process.env.FRONT_URL,
}))
