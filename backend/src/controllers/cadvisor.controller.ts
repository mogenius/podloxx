import CadvisorService from '@/services/cadvisor.service';
import { NextFunction, Request, Response } from 'express';

class CadvisorController {
  public cadvisorService = new CadvisorService();

  public getData = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
    try {
      const data: any = await this.cadvisorService.rawData();

      res.status(200).json({ data: data });
    } catch (error) {
      next(error);
    }
  };
}

export default CadvisorController;
