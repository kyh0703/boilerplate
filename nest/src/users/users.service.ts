import { HttpException, HttpStatus, Injectable, Logger } from '@nestjs/common'
import { InjectRepository } from '@nestjs/typeorm'
import { Repository, DeleteResult } from 'typeorm'
import * as argon2 from 'argon2'
import { validate } from 'class-validator'

import { User } from './entities/user.entity'
import type { CreateUserDto, UpdateUserDto } from './dto'

@Injectable()
export class UsersService {
  private logger: Logger = new Logger(UsersService.name)

  constructor(
    @InjectRepository(User)
    private usersRepository: Repository<User>,
  ) {}

  async findAll(): Promise<User[]> {
    return await this.usersRepository.find()
  }

  async findOne(email: string, password: string): Promise<User> {
    const user = await this.usersRepository.findOne({ where: { email } })
    if (!user) {
      return null
    }

    if (!(await argon2.verify(user.password, password))) {
      return null
    }

    return user
  }

  async findById(id: number): Promise<User | null> {
    return await this.usersRepository.findOneBy({ id })
  }

  async findByEmail(email: string): Promise<User | null> {
    return await this.usersRepository.findOneBy({ email })
  }

  async create(createUserDto: CreateUserDto): Promise<User> {
    const { email, username, password, photo, provider } = createUserDto
    const user = await this.usersRepository.findOneBy({ username, email })
    if (user) {
      const errors = { username: 'Username and email must be unique.' }
      throw new HttpException(
        { message: 'Input data validation failed', errors },
        HttpStatus.BAD_REQUEST,
      )
    }

    let hashPassword = null
    if (password) {
      hashPassword = await argon2.hash(password)
    }

    // create new user
    const newUser = new User()
    newUser.username = username
    newUser.email = email
    newUser.password = hashPassword
    newUser.photo = photo
    newUser.provider = provider

    const errors = await validate(newUser)
    this.logger.debug(errors)
    if (errors.length > 0) {
      // const _error = { username: 'Userinput is not valid.' }
      throw new HttpException(
        { message: 'Input data validation failed', errors },
        HttpStatus.BAD_REQUEST,
      )
    }
    const saveUser = this.usersRepository.save(newUser)
    return saveUser
  }

  async update(id: number, updateUserDto: UpdateUserDto): Promise<User> {
    const toUpdate = await this.usersRepository.findOneBy({ id })
    delete toUpdate.password

    const updated = Object.assign(toUpdate, updateUserDto)
    return await this.usersRepository.save(updated)
  }

  async delete(email: string): Promise<DeleteResult> {
    return await this.usersRepository.delete({ email })
  }

  async setCurrentRefreshToken(refreshToken: string, id: number) {
    const localRefreshToken = await argon2.hash(refreshToken)
    await this.usersRepository.update(id, { localRefreshToken })
  }

  async getRefreshToken(refreshToken: string, id: number) {
    const user = await this.findById(id)

    if (!(await argon2.verify(user.localRefreshToken, refreshToken))) {
      return null
    }

    return user
  }

  async removeRefreshToken(id: number) {
    return this.usersRepository.update(id, {
      localRefreshToken: null,
    })
  }
}
