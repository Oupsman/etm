export interface Task {
  ID: number,
  name: string,
  comment: string,
  link?: string,
  duedate: string,
  isbacklog: boolean,
  iscompleted: boolean,
  isstarted: boolean,
  priority: boolean,
  urgency: boolean,
  categoryid: number,
  CreatedAt?: string,
  UpdatedAt?: string,
}

export interface NewTask {
  name: string,
  comment: string,
  link?: string,
  duedate: string,
  isbacklog: boolean,
  categoryid: number,
}
