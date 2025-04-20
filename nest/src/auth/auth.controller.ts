import { Controller, Logger, Post, Req, Res, UseGuards } from '@nestjs/common'
import { ApiBearerAuth, ApiResponse, ApiTags } from '@nestjs/swagger'
import type { AuthService } from './auth.service'
import type { ConfigService } from '@nestjs/config'
import { AuthGuard } from '@nestjs/passport'
import type { OAuthDto } from './dto/oauth.dto'
import { JwtRefreshGuard } from './guards/jwt-refresh.guard'

@ApiTags('auth')
@Controller('auth')
export class AuthController {
  private logger: Logger = new Logger(AuthController.name)

  constructor(
    private readonly authService: AuthService,
    private readonly configService: ConfigService,
  ) {}

  @Post('refresh')
  @ApiBearerAuth()
  @UseGuards(JwtRefreshGuard)
  async refreshToken(@Req() req) {}

  @Post('google')
  @UseGuards(AuthGuard('google'))
  @ApiResponse({ status: 302, description: 'Redirect to the frontend' })
  async google(@Req() oauthDto: OAuthDto, @Res({ passthrough: true }) res) {
    this.logger.log('Google OAuth initiated', oauthDto, res)
  }

  @Post('github')
  async github(@Req() oauthDto: OAuthDto, @Res({ passthrough: true }) res) {
    this.logger.log('GitHub OAuth initiated', oauthDto, res)
  }

  @Post('naver')
  async naver(@Req() oauthDto: OAuthDto, @Res({ passthrough: true }) res) {
    this.logger.log('Naver OAuth initiated', oauthDto, res)
  }

  @Post('kakao')
  async kakao(@Req() oauthDto: OAuthDto, @Res({ passthrough: true }) res) {
    this.logger.log('Kakao OAuth initiated', oauthDto, res)
  }
}
