import { IsString } from 'class-validator'

export class RefreshDto {
  @IsString()
  readonly accessToken: string
}
