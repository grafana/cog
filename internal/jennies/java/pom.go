package java

import (
	"fmt"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/common"
)

type Pom struct {
	config Config
}

func (jenny Pom) JennyName() string {
	return "Pom"
}

func (jenny Pom) Generate(_ common.Context) (codejen.Files, error) {
	return codejen.Files{
		*codejen.NewFile("pom.xml", []byte(jenny.generatePom()), jenny),
	}, nil
}

func (jenny Pom) generatePom() string {
	return fmt.Sprintf(`
<project xmlns="http://maven.apache.org/POM/4.0.0" 
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
	xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/maven-v4_0_0.xsd">
	<modelVersion>4.0.0</modelVersion>
	<groupId>com.grafana.foundation</groupId>
	<artifactId>sdk</artifactId>
	<version>%s</version>
	<packaging>pom</packaging>

	<name>Grafana foundation SDK</name>
	<description>Grafana Java library</description>

	<url>https://github.com/grafana/grafana-foundation-sdk</url>
	<licenses>
		<license>
			<name>The Apache License, Version 2.0</name>
			<url>https://www.apache.org/licenses/LICENSE-2.0.txt</url>
		</license>
	</licenses>

	<developers>
		<developer>
			<name>Grafana</name>
			<email>platform-cat@grafana.com</email>
			<organization>Grafana</organization>
			<organizationUrl>https://grafana.com</organizationUrl>
		</developer>
	</developers>

	<properties>
		<project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
		<maven.compiler.source>11</maven.compiler.source>
		<maven.compiler.target>11</maven.compiler.target>
	</properties>

	<dependencies>
		<dependency>
			<groupId>com.fasterxml.jackson.core</groupId>
			<artifactId>jackson-databind</artifactId>
			<version>2.17.1</version>
		</dependency>
	</dependencies>
</project>
`, jenny.config.MavenVersion)
}
