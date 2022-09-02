import { Router } from 'express';
import { Routes } from '@interfaces/routes.interface';
import CadvisorController from '@/controllers/cadvisor.controller';

class CadvisorRoute implements Routes {
  public path = '/cadvisor';
  public router = Router();
  public cadvisorController = new CadvisorController();

  constructor() {
    this.initializeRoutes();
  }

  private initializeRoutes() {
    this.router.get(`${this.path}`, this.cadvisorController.getData);
    // this.router.get(`${this.path}/:id(\\d+)`, this.usersController.getUserById);
    // this.router.post(`${this.path}`, validationMiddleware(CreateUserDto, 'body'), this.usersController.createUser);
    // this.router.put(`${this.path}/:id(\\d+)`, validationMiddleware(CreateUserDto, 'body', true), this.usersController.updateUser);
    // this.router.delete(`${this.path}/:id(\\d+)`, this.usersController.deleteUser);
  }
}

export default CadvisorRoute;
