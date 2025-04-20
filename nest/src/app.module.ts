import { Module } from '@nestjs/common'
import { ConfigModule, type ConfigService } from '@nestjs/config'
import { validate } from './config/env.validation'
import { UsersModule } from './users/users.module'
import appConfig from './config/app.config'
import databaseConfig from './config/database.config'
import { AuthModule } from './auth/auth.module'
import authConfig from './config/auth.config'
import { PrismaModule } from './prisma/prisma.module'

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      load: [appConfig, databaseConfig, authConfig],
      ignoreEnvFile: process.env.NODE_ENV === 'production',
      validate,
      validationOptions: {
        abortEarly: true,
      },
    }),
    PrismaModule,
    UsersModule,
    AuthModule,
  ],
  controllers: [],
  providers: [],
})
export class AppModule {
  static port: string

  constructor(private configService: ConfigService) {
    AppModule.port = this.configService.get<string>('app.port')
  }
}
