import { IsNotEmpty, IsString } from 'class-validator'

export class OAuthDto {
  @IsNotEmpty()
  @IsString()
  readonly code: string
}
