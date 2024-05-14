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
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/maven-v4_0_0.xsd">
	<modelVersion>4.0.0</modelVersion>
	<groupId>com.grafana.foundation</groupId>
	<artifactId>foundation</artifactId>
	<version>%s</version>

	<name>${project.groupId}:${project.artifactId}</name>
	<description>Grafana Java library</description>
	<packaging>jar</packaging>
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

	<scm>
		<connection>scm:git:git@github.com:grafana/grafana-foundation-sdk.git</connection>
		<developerConnection>scm:git:ssh://github.com:grafana/grafana-foundation-sdk.git</developerConnection>
		<url>https://github.com/grafana/grafana-foundation-sdk/tree/master</url>
	</scm>

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

	<build>
		<plugin>
			<artifactId>maven-clean-plugin</artifactId>
			<version>3.1.0</version>
		</plugin>
		<plugin>
			<artifactId>maven-deploy-plugin</artifactId>
			<version>2.8.2</version>
		</plugin>
	</build>

</project>
`, jenny.config.MavenVersion)
}
