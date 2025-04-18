import { Request } from 'express'

export interface OAuthUser {
  provider: string
  email: string
  username: string
  photo?: string
}

export type RequestWithSocialUser = Request & { user: OAuthUser }
