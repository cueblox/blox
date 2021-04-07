{
    // No "version", we expect people to use the path
    // of the schema to version
    _schema: {
        name: "contact"
        namespace: "schemas.cueblox.com"
    }

    #Contact: {
        _model: {
            // Lets assume lowercase Profile is ID
            // Lets assume lowercase Profile with _id is the foreign key
            // Plural for directory name
            plural: "contacts"
        }

        name: #Name
        address: #Address
        phone: string
        email: string
    }

    #Name: {
        forename: string
        surname: string
    }

    #Address: {
        number: string
        street: string
        city: string
        country: string
        postcode: string
    }
}
