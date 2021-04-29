import SwaggerUI from "swagger-ui-react"
import "swagger-ui-react/swagger-ui.css"

export default function Swagger() {
    return <SwaggerUI url={`${process.env.NEXT_PUBLIC_API_BASE}/swagger.yaml`} />
}
