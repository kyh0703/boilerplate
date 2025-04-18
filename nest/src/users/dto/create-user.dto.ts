import { IsString, IsEmail, IsNotEmpty, IsOptional } from 'class-validator'

export class CreateUserDto {
  @IsNotEmpty()
  @IsString()
  readonly provider: string

  @IsNotEmpty()
  @IsEmail()
  readonly email: string

  @IsNotEmpty()
  @IsString()
  readonly username: string

  @IsString()
  readonly password: string

  @IsOptional()
  @IsString()
  readonly photo?: string
}
