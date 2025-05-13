export interface Task {
  ID:          number,
  name:        string,
  comment:     string,
  duedate:     string,
  isbacklog:   boolean,
  iscompleted: boolean,
  priority:    boolean,
  urgency:     boolean,
  categoryid:  number,
}
