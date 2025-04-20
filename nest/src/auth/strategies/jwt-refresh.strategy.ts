import { Injectable, UnauthorizedException } from '@nestjs/common'
import { PassportStrategy } from '@nestjs/passport'
import { ExtractJwt, Strategy } from 'passport-jwt'
import { ConfigService } from '@nestjs/config'
import { Request } from 'express'
import type { PrismaService } from 'src/prisma/prisma.service'

@Injectable()
export class JwtRefreshStrategy extends PassportStrategy(
  Strategy,
  'jwt-refresh'
) {
  constructor(
    private configService: ConfigService,
    private prisma: PrismaService
  ) {
    super({
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      secretOrKey: configService.get<string>('auth.refreshTokenSecret'),
      passReqToCallback: true,
    })
  }

  async validate(req: Request, payload: any) {
    const refreshToken = req.headers.authorization.replace('Bearer ', '').trim()

    // Check if the token exists in the database
    const tokenRecord = await this.prisma.token.findFirst({
      where: {
        userId: payload.sub,
        token: refreshToken,
        expiresAt: {
          gt: new Date(),
        },
      },
    })

    if (!tokenRecord) {
      throw new UnauthorizedException('Invalid refresh token')
    }

    return {
      ...payload,
      refreshToken,
    }
  }
}
