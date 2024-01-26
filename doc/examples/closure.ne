obj DatabaseRecord {
    id: int
    name: string
    data: string

    get id, name, data
}

let database = [DatabaseRecord]

let createDBController = () => {
    let lastID = 0

    let validateRecord = (record: DatabaseRecord) => bool {
        if record.name is nil || record.data is nil {
            => false
        }
        => true
    }

    let insertRecord = (record: DatabaseRecord) => string {
        if validateRecord record {
            lastID = lastID + 1
            record.id = lastID
            database.append record
            => "Inserted successfully with ID: {}" lastID
        }
        => "Failed to insert record."
    }

    let readRecord = (id: int) => DatabaseRecord? {
        for record in database {
            if record.id is id {
                => record
            }
        }
        => nil  // Record not found.
    }

    let deleteRecord = (id: int) => string {
        for index, record in enumerate(database) {
            if record.id is id {
                database = database[:index] + database[index+1:]
                => "Record with ID: {} deleted successfully." id
            }
        }
        => "Record not found."
    }

    => |
        "insert": insert_record, 
        "read": read_record, 
        "delete": delete_record
    |
}

let dbController = createDBController

let newRecord = DatabaseRecord{id: 0, name: "test", data: "sample data"}
println dbController["insert"] newRecord

let fetchedRecord = dbController["read"] 1
println "Fetched Record Name: {}, Data: {}" fetchedRecord?.name fetchedRecord?.data

println dbController["delete"] 1
