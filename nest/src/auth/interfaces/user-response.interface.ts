import { Request } from 'express'

export interface UserRO {
  id: number
  email: string
  username: string
  photo?: string
}

export type RequestWithUser = Request & { user: UserRO }
