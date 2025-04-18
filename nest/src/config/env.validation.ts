import { plainToInstance } from 'class-transformer'
import {
  IsEnum,
  validateSync,
  IsNotEmpty,
  IsString,
  IsNumber,
} from 'class-validator'

enum Environment {
  Development = 'development',
  Production = 'production',
  Test = 'test',
}

class EnvironmentVariables {
  @IsEnum(Environment)
  NODE_ENV: Environment

  @IsNotEmpty()
  @IsString()
  DB_TYPE: string

  @IsNotEmpty()
  @IsString()
  DB_DATABASE: string

  @IsNotEmpty()
  @IsString()
  JWT_ACCESS_TOKEN_SECRET: string

  @IsNotEmpty()
  @IsNumber()
  JWT_ACCESS_TOKEN_EXPIRES_IN: number

  @IsNotEmpty()
  @IsString()
  JWT_REFRESH_TOKEN_SECRET: string

  @IsNotEmpty()
  @IsNumber()
  JWT_REFRESH_TOKEN_EXPIRES_IN: number

  @IsString()
  GOOGLE_CLIENT_ID: string

  @IsString()
  GOOGLE_CLIENT_SECRET: string

  @IsString()
  GOOGLE_REDIRECT_URI: string

  @IsString()
  GITHUB_CLIENT_ID: string

  @IsString()
  GITHUB_CLIENT_SECRET: string

  @IsString()
  GITHUB_REDIRECT_URI: string

  @IsString()
  NAVER_CLIENT_ID: string

  @IsString()
  NAVER_CLIENT_SECRET: string

  @IsString()
  NAVER_REDIRECT_URI: string

  @IsString()
  KAKAO_CLIENT_ID: string

  @IsString()
  KAKAO_CLIENT_SECRET: string

  @IsString()
  KAKAO_REDIRECT_URI: string
}

export function validate(config: Record<string, unknown>) {
  const validateConfig = plainToInstance(EnvironmentVariables, config, {
    enableImplicitConversion: true,
  })
  const errors = validateSync(validateConfig, { skipMissingProperties: false })

  if (errors.length > 0) {
    throw new Error(errors.toString())
  }

  return validateConfig
}
