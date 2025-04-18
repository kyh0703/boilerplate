import { Module } from '@nestjs/common'
import { ConfigModule } from '@nestjs/config'
import { AppController } from './app.controller'
import { AppService } from './app.service'
import { AuthController } from './auth/auth.controller'
import { AuthService } from './auth/auth.service'
import { validate } from './config/env.validation'
import { OauthModule } from './oauth/oauth.module'
import { UsersController } from './users/users.controller'
import { UsersModule } from './users/users.module'

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      ignoreEnvFile: process.env.NODE_ENV === 'production',
      validate,
    }),
    OauthModule,
    UsersModule,
  ],
  controllers: [AppController, AuthController, UsersController],
  providers: [AppService, AuthService],
})
export class AppModule {}
