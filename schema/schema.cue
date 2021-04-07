{
    // No "version", we expect people to use the path
    // of the schema to version
    _schema: {
        name: "blox"
        namespace: "schemas.devrel-blox.com"
    }

    #Profile: {
        _model: {
            // Lets assume lowercase Profile is ID
            // Lets assume lowercase Profile with _id is the foreign key
            // Plural for directory name
            plural: "profiles"
        }

        name: #Name
        address: #Address
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
