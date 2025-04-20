import { registerAs } from '@nestjs/config'

export default registerAs('auth', () => ({
  accessTokenSecret: process.env.JWT_ACCESS_TOKEN_SECRET,
  accessTokenExpiresIn: process.env.JWT_ACCESS_TOKEN_EXPIRES_IN,
  refreshTokenSecret: process.env.JWT_REFRESH_TOKEN_SECRET,
  refreshTokenExpiresIn: process.env.JWT_REFRESH_TOKEN_EXPIRES_IN,
  googleId: process.env.GOOGLE_CLIENT_ID,
  googleSecret: process.env.GOOGLE_CLIENT_SECRET,
  googleRedirectURI: process.env.GOOGLE_REDIRECT_URI,
  githubId: process.env.GITHUB_CLIENT_ID,
  githubSecret: process.env.GITHUB_CLIENT_SECRET,
  githubRedirectURI: process.env.GITHUB_REDIRECT_URI,
  naverId: process.env.NAVER_CLIENT_ID,
  naverSecret: process.env.NAVER_CLIENT_SECRET,
  naverRedirectURI: process.env.NAVER_REDIRECT_URI,
  kakaoId: process.env.KAKAO_CLIENT_ID,
  kakaoSecret: process.env.KAKAO_CLIENT_SECRET,
  kakaoRedirectURI: process.env.KAKAO_REDIRECT_URI,
}))
