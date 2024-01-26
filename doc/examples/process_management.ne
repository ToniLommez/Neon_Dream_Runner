obj Process {
    pid: int
    name: string
    status: string? = "Waiting"

    get pid, name, status
    set status

    fn to_string = "PID: {}, Name: {}, Status: {}" pid name status
}

obj ProcessManager {
    process_list: [Process] = []

    fn add_process(name: string) {
        let new_PID = (len process_list) + 1
        let new_Process = Process{pid: new_PID, name: name}
        process_list.append new_Process
    }

    fn terminate_process(pid: int) => bool {
        for process in process_list {
            if process.pid is pid {
                process.status = "Terminated"
                => true  // Successfully terminated
            }
        }
        => false  // PID not found
    }

    fn run_process(pid: int) => bool {
        for process in process_list {
            if process.pid is pid and process.status not in ["Terminated", "Running"] {
                process.status = "Running"
                => true  // Successfully started
            }
        }
        => false  // PID not found or process is already running/terminated
    }

    fn show_all_processes {
        for p in process_list {
            println p
        }
    }
}

fn main {
    let pm = ProcessManager{}

    pm.add_process "TextEditor"
    pm.add_process "WebBrowser"
    pm.add_process "FileExplorer"

    pm.run_process 2
    pm.terminate_process 3

    pm.show_all_processes
}