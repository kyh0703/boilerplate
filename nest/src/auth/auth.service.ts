import {
  BadRequestException,
  Injectable,
  Logger,
  UnauthorizedException,
} from '@nestjs/common'
import { ConfigService } from '@nestjs/config'
import { JwtService } from '@nestjs/jwt'
import { Response } from 'express'
import type { CreateUserDto } from '../users/dto'
import type { UsersService } from '../users/users.service'
import { Payload, Token, UserRO } from './interfaces'

@Injectable()
export class AuthService {
  private logger: Logger = new Logger(AuthService.name)

  constructor(
    private oauthFactory: OAuthFactory,
    private usersService: UsersService,
    private jwtService: JwtService,
    private configService: ConfigService,
  ) {}

  async validateUser(email: string, pass: string): Promise<UserRO> {
    const user = await this.usersService.findByEmail(email)
    if (!user) {
      throw new BadRequestException('존재하지 않는 이메일입니다')
    }

    if (user.signinFailCount >= 5) {
      throw new UnauthorizedException('비밀번호를 5회 이상 틀렸습니다')
    }

    const isPasswordValidated: boolean = await user.validatePassword(pass)
    if (!isPasswordValidated) {
      const signinFailCount = user.signinFailCount
      this.usersService.update(user.id, {
        signinFailCount: signinFailCount + 1,
      })
      throw new UnauthorizedException(
        `비밀번호가 올바르지 않습니다. 로그인 실패 횟수 초과 시 계정이 잠금 처리됩니다 (${signinFailCount}/5)`,
      )
    }

    const validateUser: UserRO = {
      id: user.id,
      email: user.email,
      username: user.username,
      photo: user.photo,
    }

    return validateUser
  }

  async signin(signinDto: SigninDto, res: Response): Promise<UserRO> {
    const { id, email, username } = signinDto

    const accessToken = this.genAccessToken({ sub: id, email })
    const refreshToken = this.genRefreshToken({
      sub: id,
      email: email,
    })

    this.usersService.setCurrentRefreshToken(refreshToken, id)

    res.cookie('accessToken', accessToken, {
      expires: this.getAccessTokenExpiresIn(),
      httpOnly: true,
    })

    res.cookie('refreshToken', refreshToken, {
      expires: this.getRefreshTokenExpiresIn(),
      httpOnly: true,
    })

    const userRO: UserRO = {
      id,
      email,
      username,
    }

    return userRO
  }

  async signinWithOAuth2(
    code: string,
    provider: Provider,
    res: Response,
  ): Promise<Token> {
    this.logger.debug('siginWithOAuth2: ')
    const oauther = this.oauthFactory.NewStrategy(provider)
    const oauthToken = await oauther.getToken(code)
    this.logger.debug('emit oauth token')
    const user = await oauther.getUser(oauthToken.access_token)
    this.logger.debug('get user info')

    const saveUser = await this.usersService.findByEmail(user.email)
    if (saveUser && saveUser.provider !== user.provider) {
      throw new BadRequestException('현재 계정으로 가입 한 메일이 존재합니다')
    }

    this.logger.debug(user)
    let createUser = saveUser
    if (!saveUser) {
      createUser = await this.usersService.create(user)
    }
    this.logger.debug('saveUser: ${saveUser}')

    const accessToken = this.genAccessToken({
      sub: createUser.id,
      email: createUser.email,
    })

    const refreshToken = this.genRefreshToken({
      sub: createUser.id,
      email: createUser.email,
    })

    res.cookie('refreshToken', refreshToken, {
      expires: this.getRefreshTokenExpiresIn(),
      httpOnly: true,
    })

    const token: Token = {
      accessToken: accessToken,
      expiresIn: this.getAccessTokenExpiresIn().getTime(),
    }

    this.logger.debug(token)
    return token
  }

  async signup(signupDto: SignupDto): Promise<void> {
    this.logger.debug(`signup: ${signupDto}`)

    const _user = await this.usersService.findByEmail(signupDto.email)
    if (_user) {
      throw new BadRequestException('이미 존재하는 이메일입니다.')
    }

    const createUserDto: CreateUserDto = {
      provider: Provider.Local,
      email: signupDto.email,
      password: signupDto.password,
      username: signupDto.username,
    }

    const user = await this.usersService.create(createUserDto)
    if (!user) {
      throw new BadRequestException('회원가입에 실패했습니다.')
    }

    return null
  }

  async slientRefresh(user: UserRO, res: Response): Promise<UserRO> {
    const accessToken = this.genAccessToken({ sub: user.id, email: user.email })
    const refreshToken = this.genRefreshToken({
      sub: user.id,
      email: user.email,
    })

    res.cookie('accessToken', accessToken, {
      expires: this.getAccessTokenExpiresIn(),
      httpOnly: true,
    })

    res.cookie('refreshToken', refreshToken, {
      expires: this.getRefreshTokenExpiresIn(),
      httpOnly: true,
    })

    return user
  }

  private genAccessToken(payload: Payload): string {
    const accessSecret = this.configService.get('jwt.accessTokenSecret')
    const accessExpiresIn = this.configService.get('jwt.accessTokenExpiresIn')

    return this.jwtService.sign(payload, {
      secret: accessSecret,
      expiresIn: `${accessExpiresIn}h`,
    })
  }

  private genRefreshToken(payload: Payload): string {
    const refreshSecret = this.configService.get('jwt.refreshTokenSecret')
    const refreshExpiresIn = this.configService.get('jwt.refreshTokenExpiresIn')

    return this.jwtService.sign(payload, {
      secret: refreshSecret,
      expiresIn: `${refreshExpiresIn} days`,
    })
  }

  private getAccessTokenExpiresIn(): Date {
    const expires = new Date()
    const accessTokenExpiresIn = this.configService.get(
      'jwt.accessTokenExpiresIn',
    )
    expires.setDate(expires.getHours() + +accessTokenExpiresIn)
    return expires
  }

  private getRefreshTokenExpiresIn(): Date {
    const expires = new Date()
    const refreshExpiresIn = this.configService.get('jwt.refreshTokenExpiresIn')
    expires.setDate(expires.getDate() + +refreshExpiresIn)
    return expires
  }
}
