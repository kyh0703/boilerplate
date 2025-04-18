import { Injectable, UnauthorizedException } from '@nestjs/common'
import { ConfigService } from '@nestjs/config'
import axios, { AxiosResponse } from 'axios'
import * as qs from 'qs'
import { Provider } from 'src/modules/users/entities/user.entity'
import { OAuther, OAuthToken, OAuthUser } from '../interfaces'

@Injectable()
export class KakaoOAuthStrategy implements OAuther {
  constructor(protected configService: ConfigService) {}

  async getToken(code: string): Promise<OAuthToken> {
    try {
      const response: AxiosResponse = await axios({
        url: 'https://kauth.kakao.com/oauth/token',
        method: 'POST',
        headers: {
          'content-type': 'application/x-www-form-urlencoded',
        },
        data: qs.stringify({
          grant_type: 'authorization_code',
          client_id: this.configService.get('oauth.kakaoId'),
          client_secret: this.configService.get('oauth.kakaoSecret'),
          redirect_uri: this.configService.get('oauth.kakaoCallback'),
          code,
        }),
      })
      return response.data
    } catch (error) {
      throw new UnauthorizedException('카카오 인증을 실패하였습니다')
    }
  }

  async getUser(token: string): Promise<OAuthUser> {
    try {
      const response: AxiosResponse = await axios({
        url: 'https://kapi.kakao.com/v2/user/me',
        method: 'GET',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      console.log(response.data)
      const { id, kakao_account, properties } = response.data
      const { email } = kakao_account
      const { nickname } = properties
      const user: OAuthUser = {
        provider: Provider.Kakao,
        email,
        password: '' + id,
        username: nickname,
      }
      return user
    } catch (error) {
      throw new UnauthorizedException('카카오 인증을 실패하였습니다')
    }
  }
}
