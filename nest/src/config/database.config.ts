import { registerAs } from '@nestjs/config'

export default registerAs('db', () => ({
  type: process.env.DB_TYPE,
  database: process.env.DB_DATABASE,
  host: process.env.DB_HOST || 5432,
  port: parseInt(process.env.DB_PORT, 10) || 6000,
  username: process.env.DB_USERNAME,
  password: process.env.DB_PASSWORD,
  synchronize: process.env.DB_SYNCHRONIZE === 'true',
}))
