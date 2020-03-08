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

const getLists = async () => {
  const res = await fetch(`${API_URL}/list/all`)
  return res.json() as Promise<ListSimple[]>
}

const createList = async (name: string) => {
  const res = await fetch(`${API_URL}/list/new?name=`+encodeURIComponent(name), {
    method: 'post',
  })
  return res.json() as Promise<ListSimple>
}

const deleteList = async (listID: string) => {
  const res = await fetch(`${API_URL}/list/id/${listID}?return_list=true`, {
    method: 'delete',
  })
  return res.json() as Promise<ListSimple[]>
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

  const handleChange = function(event: React.ChangeEvent<HTMLInputElement>) {
    setCurrent(event.target.value)
  }

  const handleSubmit = function(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    createList(current).then((res) => {
      setCurrent("")
      props.callback(res.UUID)
    })
  }

  const handleDelete = function(id: string) {
    deleteList(id).then(res => {
      setLists(res)
    })
  }

  return (
    <div className="todoLists">
      <div className="header">
        <form onSubmit={handleSubmit}>
          <input placeholder="Title" value={current} onChange={handleChange} />
          <button type="submit"> Create new list</button>
          {lists && lists.map((value) => (
            <li key={value.UUID}>
              <a onClick={() => props.callback(value.UUID)}>{value.Name} </a>
              <a onClick={() => handleDelete(value.UUID)}>X</a>
            </li>
          ))}
        </form>
      </div>
    </div>
  )
}

const getTasks = async (id: string) => {
  const res = await fetch(`${API_URL}/list/id/${id}`)
  return res.json() as Promise<Tasks>
}

const addTask = async (listID: string, text: string) => {
  const res = await fetch(`${API_URL}/list/id/${listID}/item/add?text=${encodeURIComponent(text)}&return_list=true`, {
    method: 'post',
  })
  return res.json() as Promise<Tasks>
}

const deleteTask = async (listID: string, taskID: string) => {
  const res = await fetch(`${API_URL}/list/id/${listID}/item/id/${taskID}?return_list=true`, {
    method: 'delete',
  })
  return res.json() as Promise<Tasks>
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
    addTask(props.UUID, current).then((res) => {
      setData(res)
    })
    setCurrent("")
  }

  const handleDelete = function(listID: string, taskID: string) {
    deleteTask(props.UUID, taskID).then(res => {
      setData(res)
    })
  }

  return (
    <div className="todoListMain">
      <div className="header">
      <h1>{data.Name}</h1>
        <form onSubmit={handleSubmit}>
          <input placeholder="Task" value={current} onChange={handleChange} />
          <button type="submit"> Add Task </button>
          {data.List && data.List.Items && data.List.Items.map((value) => (
            <li key={value.UUID} onClick={() => handleDelete(props.UUID, value.UUID)}>
              {value.value}
            </li>
          ))}
        </form>
      </div>
    </div>
  )
}

export default TodoApp;
