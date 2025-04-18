import { Injectable, UnauthorizedException } from '@nestjs/common'
import { ConfigService } from '@nestjs/config'
import axios, { AxiosResponse } from 'axios'
import { Provider } from 'src/modules/users/entities/user.entity'
import { OAuther, OAuthToken, OAuthUser } from '../interfaces'

@Injectable()
export class GoogleOAuthStrategy implements OAuther {
  constructor(protected configService: ConfigService) {}

  async getToken(code: string): Promise<OAuthToken> {
    try {
      const response: AxiosResponse = await axios({
        url: 'https://oauth2.googleapis.com/token',
        method: 'POST',
        headers: {
          'Content-Type': 'x-www-form-urlencoded',
        },
        data: {
          code,
          client_id: this.configService.get('oauth.googleId'),
          client_secret: this.configService.get('oauth.googleSecret'),
          redirect_uri: this.configService.get('oauth.googleRedirectURI'),
          grant_type: 'authorization_code',
        },
      })
      return response.data
    } catch (error) {
      throw new UnauthorizedException('구글 인증을 실패하였습니다')
    }
  }

  async getUser(token: string): Promise<OAuthUser> {
    try {
      const response: AxiosResponse = await axios({
        url: 'https://www.googleapis.com/oauth2/v2/userinfo',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      console.log(response.data)
      const { id, email, avatar_url, name, bio } = response.data
      const user: OAuthUser = {
        provider: Provider.Google,
        email,
        password: '' + id,
        username: name,
        photo: avatar_url,
        bio: bio,
      }
      return user
    } catch (error) {
      throw new UnauthorizedException('구글 인증을 실패하였습니다')
    }
  }
}
