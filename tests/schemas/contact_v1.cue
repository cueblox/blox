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

        name: #ContactName
        address: #ContactAddress
        phone: string
        email: string
    }

    #ContactName: {
        forename: string
        surname: string
    }

    #ContactAddress: {
        number: string
        street: string
        city: string
        country: string
        postcode: string
    }
}
