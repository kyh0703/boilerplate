import type { OAuthToken } from './oauth-token.interface'
import { OAuthUser } from './oauth-user.interface'

export interface OAuther {
  getToken(code: string): Promise<OAuthToken>
  getUser(token: string): Promise<OAuthUser>
}
