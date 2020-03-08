import React, { FunctionComponent, useState, useEffect } from 'react'

const API_URL = 'http://localhost:8080'

interface Props {
}

interface Callback {
  (arg0: string): void,
}

interface ListProps {
  callback: Callback,
}

interface TaskProps {
  UUID: string,
}

interface ListSimple {
  UUID: string,
  Name: string,
  Created: Date,
}

interface Item {
  UUID: string;
  value: string;
}

interface List {
  Items: Item[];
  MaxTotal: number;
  MaxTextLength: number;
  Total: number;
}

interface Tasks {
  UUID: string;
  Name: string;
  Created: Date;
  List: List;
}

const getLists= async () => {
  const res = await fetch(`${API_URL}/list/all`)
  return res.json() as Promise<ListSimple[]>
}

const getTasks= async (id: string) => {
  const res = await fetch(`${API_URL}/list/id/${id}`)
  return res.json() as Promise<Tasks>
}

const TodoApp: FunctionComponent = () => {
  const [selectedList, setSelectedList] = useState("")
  const selectedListFunc = function(id: string) {
    setSelectedList(id)
  }
  return (
    <div className="todoListMain">
      { selectedList === "" ? (
        <TodoAppLists callback={selectedListFunc}/>
      ): (
        <TodoAppTasks UUID={selectedList} />
      )}
    </div>
  )
}

const TodoAppLists: FunctionComponent<ListProps> = (props) => {
  const [current, setCurrent] = useState("")
  const [lists, setLists] = useState(Array<ListSimple>())

  useEffect(() => {
    getLists().then(
      res => {
        setLists(res)
      }
    )
  },[])

  useEffect(() => {
    console.log(lists)
  },[lists])

  const handleChange = function(event: React.ChangeEvent<HTMLInputElement>) {
    setCurrent(event.target.value)
  }

  const handleSubmit = function(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    console.log("create new list:", current)
    setCurrent("")
  }

  return (
    <div className="todoLists">
      <div className="header">
        <form onSubmit={handleSubmit}>
          <input placeholder="Title" value={current} onChange={handleChange} />
          <button type="submit"> Create new list</button>
          {lists.map((value) => (
            <li key={value.UUID} onClick={() => props.callback(value.UUID)}>
              {value.Name}
            </li>
          ))}
        </form>
      </div>
    </div>
  )
}

const TodoAppTasks: FunctionComponent<TaskProps> = (props) => {
  const [data, setData] = useState({} as Tasks)
  const [current, setCurrent] = useState("")

  useEffect(() => {
    getTasks(props.UUID).then(
      res => {
        setData(res)
      }
    )
  },[])

  const handleChange = function(event: React.ChangeEvent<HTMLInputElement>) {
    setCurrent(event.target.value)
  }

  const handleSubmit = function(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    console.log("add to tasks:", current)
    setCurrent("")
  }

  return (
    <div className="todoListMain">
      <div className="header">
        <form onSubmit={handleSubmit}>
          <input placeholder="Task" value={current} onChange={handleChange} />
          <button type="submit"> Add Task </button>
          {data.List && data.List.Items.map((value) => (
            <li key={value.UUID}>
              {value.value}
            </li>
          ))}
        </form>
      </div>
    </div>
  )
}

export default TodoApp;
