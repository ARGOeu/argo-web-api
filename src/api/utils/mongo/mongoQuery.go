/*
 * Copyright (c) 2013 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the
 * License. You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an "AS
 * IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language
 * governing permissions and limitations under the License.
 *
 * The views and conclusions contained in the software and
 * documentation are those of the authors and should not be
 * interpreted as representing official policies, either expressed
 * or implied, of either GRNET S.A., SRCE or IN2P3 CNRS Computing
 * Centre
 *
 * The work represented by this source file is partially funded by
 * the EGI-InSPIRE project through the European Commission's 7th
 * Framework Programme (contract # INFSO-RI-261323)
 */

package mongo

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func openCollection(session *mgo.Session, dbName string, collectionName string) *mgo.Collection{
	
	c := session.DB(dbName).C(collectionName)
	
	return c
}

func Pipe(session *mgo.Session, dbName string, collectionName string, query []bson.M, results interface{}) error{
	
	c := openCollection(session,dbName,collectionName)
	
	err := c.Pipe(query).All(results)
	
	return err
	
}

func Find(session *mgo.Session, dbName string, collectionName string, query []bson.M, results interface{}) error{
	
	c := openCollection(session,dbName,collectionName)
	
	err := c.Find(query).All(results)
	
	return err
	
}

func Insert(session *mgo.Session, dbName string, collectionName string, query bson.M) error{
	
	c := openCollection(session,dbName,collectionName)
	
	err := c.Insert(query)
	
	return err
	
}