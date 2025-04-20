import { Injectable, UnauthorizedException } from '@nestjs/common'
import { ConfigService } from '@nestjs/config'
import axios, { AxiosResponse } from 'axios'
import { Provider } from 'src/modules/users/entities/user.entity'
import { OAuther, OAuthToken, OAuthUser } from '../interfaces'

@Injectable()
export class GithubStrategy implements OAuther {
  constructor(protected configService: ConfigService) {}

  async getToken(code: string): Promise<OAuthToken> {
    try {
      const response: AxiosResponse = await axios({
        url: 'https://github.com/login/oauth/access_token',
        method: 'POST',
        headers: {
          accept: 'application/json',
        },
        data: {
          code,
          client_id: this.configService.get('oauth.githubId'),
          client_secret: this.configService.get('oauth.githubSecret'),
        },
      })
      return response.data
    } catch (error) {
      throw new UnauthorizedException('깃허브 인증을 실패하였습니다')
    }
  }

  async getUser(token: string): Promise<OAuthUser> {
    try {
      const response: AxiosResponse = await axios({
        url: 'https://api.github.com/user',
        headers: {
          Authorization: `token ${token}`,
        },
      })
      const { id, email, avatar_url, name, bio } = response.data
      const user: OAuthUser = {
        provider: Provider.Github,
        email,
        password: '' + id,
        username: name,
        photo: avatar_url,
        bio,
      }
      return user
    } catch (error) {
      throw new UnauthorizedException('깃허브 인증을 실패하였습니다')
    }
  }
}
