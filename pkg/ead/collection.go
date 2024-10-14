package ead

import (
	"github.com/lestrrat-go/libxml2/types"
)

type Collection struct {
	SolrAddMessage string
	XPathParts     CollectionXPathParts
}

type CollectionXPathParts struct {
	Abstract           CollectionXPathPart
	Acqinfo            CollectionXPathPart
	Appraisal          CollectionXPathPart
	Author             CollectionXPathPart
	Bioghist           CollectionXPathPart
	Chronlist          CollectionXPathPart
	Collection         CollectionXPathPart
	Corpname           CollectionXPathPart
	Creator            CollectionXPathPart
	Custodhist         CollectionXPathPart
	Ead                CollectionXPathPart
	Famname            CollectionXPathPart
	Function           CollectionXPathPart
	Genreform          CollectionXPathPart
	Geogname           CollectionXPathPart
	Heading            CollectionXPathPart
	Language           CollectionXPathPart
	Material_type      CollectionXPathPart
	Name               CollectionXPathPart
	Note               CollectionXPathPart
	Occupation         CollectionXPathPart
	Persname           CollectionXPathPart
	Phystech           CollectionXPathPart
	Place              CollectionXPathPart
	Scopecontent       CollectionXPathPart
	Subject            CollectionXPathPart
	Subject_facets     CollectionXPathPart
	Title              CollectionXPathPart
	Unitdate           CollectionXPathPart
	Unitdate_bulk      CollectionXPathPart
	Unitdate_end       CollectionXPathPart
	Unitdate_inclusive CollectionXPathPart
	Unitdate_normal    CollectionXPathPart
	Unitdate_start     CollectionXPathPart
	Unitid             CollectionXPathPart
	Unittitle          CollectionXPathPart
}

type CollectionXPathPart struct {
	Query  string
	Values []string
}

func (collection *Collection) populateXPathParts(node types.Node) error {
	var err error

	xp := &collection.XPathParts

	xp.Abstract.Query = "//archdesc[@level='collection']/did/abstract"
	xp.Abstract.Values, err = parseValues(xp.Abstract.Query, node)
	if err != nil {
		return err
	}

	xp.Acqinfo.Query = "//archdesc[@level='collection']/acqinfo/p"
	xp.Acqinfo.Values, err = parseValues(xp.Acqinfo.Query, node)
	if err != nil {
		return err
	}

	xp.Appraisal.Query = "//archdesc[@level='collection']/appraisal/p"
	xp.Appraisal.Values, err = parseValues(xp.Appraisal.Query, node)
	if err != nil {
		return err
	}

	xp.Author.Query = "//filedesc/titlestmt/author"
	xp.Author.Values, err = parseValues(xp.Author.Query, node)
	if err != nil {
		return err
	}

	xp.Bioghist.Query = "//archdesc[@level='collection']/bioghist/p"
	xp.Bioghist.Values, err = parseValues(xp.Bioghist.Query, node)
	if err != nil {
		return err
	}

	xp.Chronlist.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//chronlist/chronitem//text()"
	xp.Chronlist.Values, err = parseValues(xp.Chronlist.Query, node)
	if err != nil {
		return err
	}

	xp.Collection.Query = "//archdesc[@level='collection']/did/unittitle"
	xp.Collection.Values, err = parseValues(xp.Collection.Query, node)
	if err != nil {
		return err
	}

	xp.Corpname.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//corpname"
	xp.Corpname.Values, err = parseValues(xp.Corpname.Query, node)
	if err != nil {
		return err
	}

	xp.Creator.Query = "//archdesc[@level='collection']/did/origination[@label='creator']/*[#{creator_fields_to_xpath}]"
	xp.Creator.Values, err = parseValues(xp.Creator.Query, node)
	if err != nil {
		return err
	}

	xp.Custodhist.Query = "//archdesc[@level='collection']/custodhist/p"
	xp.Custodhist.Values, err = parseValues(xp.Custodhist.Query, node)
	if err != nil {
		return err
	}

	xp.Ead.Query = "//eadid"
	xp.Ead.Values, err = parseValues(xp.Ead.Query, node)
	if err != nil {
		return err
	}

	xp.Famname.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//famname"
	xp.Famname.Values, err = parseValues(xp.Famname.Query, node)
	if err != nil {
		return err
	}

	xp.Function.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//function"
	xp.Function.Values, err = parseValues(xp.Function.Query, node)
	if err != nil {
		return err
	}

	xp.Genreform.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//genreform"
	xp.Genreform.Values, err = parseValues(xp.Genreform.Query, node)
	if err != nil {
		return err
	}

	xp.Geogname.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//geogname"
	xp.Geogname.Values, err = parseValues(xp.Geogname.Query, node)
	if err != nil {
		return err
	}

	xp.Heading.Query = "//archdesc[@level='collection']/did/unittitle"
	xp.Heading.Values, err = parseValues(xp.Heading.Query, node)
	if err != nil {
		return err
	}

	xp.Language.Query = "//archdesc[@level='collection']/did/langmaterial/language/@langcode"
	xp.Language.Values, err = parseValues(xp.Language.Query, node)
	if err != nil {
		return err
	}

	xp.Material_type.Query = "//genreform"
	xp.Material_type.Values, err = parseValues(xp.Material_type.Query, node)
	if err != nil {
		return err
	}

	xp.Name.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//name"
	xp.Name.Values, err = parseValues(xp.Name.Query, node)
	if err != nil {
		return err
	}

	xp.Note.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//note"
	xp.Note.Values, err = parseValues(xp.Note.Query, node)
	if err != nil {
		return err
	}

	xp.Occupation.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//occupation"
	xp.Occupation.Values, err = parseValues(xp.Occupation.Query, node)
	if err != nil {
		return err
	}

	xp.Persname.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//persname"
	xp.Persname.Values, err = parseValues(xp.Persname.Query, node)
	if err != nil {
		return err
	}

	xp.Phystech.Query = "//archdesc[@level='collection']/phystech/p"
	xp.Phystech.Values, err = parseValues(xp.Phystech.Query, node)
	if err != nil {
		return err
	}

	xp.Place.Query = "//geogname"
	xp.Place.Values, err = parseValues(xp.Place.Query, node)
	if err != nil {
		return err
	}

	xp.Scopecontent.Query = "//archdesc[@level='collection']/scopecontent/p"
	xp.Scopecontent.Values, err = parseValues(xp.Scopecontent.Query, node)
	if err != nil {
		return err
	}

	xp.Subject.Query = "//[archdesc[@level='collection']/*[name() != 'dsc']//subject"
	xp.Subject.Values, err = parseValues(xp.Subject.Query, node)
	if err != nil {
		return err
	}

	xp.Subject_facets.Query = "//*[local-name()='subject' or local-name()='function' or local-name() = 'occupation']"
	xp.Subject_facets.Values, err = parseValues(xp.Subject_facets.Query, node)
	if err != nil {
		return err
	}

	xp.Title.Query = "//archdesc[@level='collection']/*[name() != 'dsc']//title"
	xp.Title.Values, err = parseValues(xp.Title.Query, node)
	if err != nil {
		return err
	}

	xp.Unitdate.Query = "//archdesc[@level='collection']/did/unitdate[not(@type)]"
	xp.Unitdate.Values, err = parseValues(xp.Unitdate.Query, node)
	if err != nil {
		return err
	}

	xp.Unitdate_bulk.Query = "//archdesc[@level='collection']/did/unitdate[@type='bulk']"
	xp.Unitdate_bulk.Values, err = parseValues(xp.Unitdate_bulk.Query, node)
	if err != nil {
		return err
	}

	xp.Unitdate_end.Query = "//archdesc[@level='collection']/did/unitdate/@normal"
	xp.Unitdate_end.Values, err = parseValues(xp.Unitdate_end.Query, node)
	if err != nil {
		return err
	}

	xp.Unitdate_inclusive.Query = "//archdesc[@level='collection']/did/unitdate[@type='inclusive']"
	xp.Unitdate_inclusive.Values, err = parseValues(xp.Unitdate_inclusive.Query, node)
	if err != nil {
		return err
	}

	xp.Unitdate_normal.Query = "//archdesc[@level='collection']/did/unitdate/@normal"
	xp.Unitdate_normal.Values, err = parseValues(xp.Unitdate_normal.Query, node)
	if err != nil {
		return err
	}

	xp.Unitdate_start.Query = "//archdesc[@level='collection']/did/unitdate/@normal"
	xp.Unitdate_start.Values, err = parseValues(xp.Unitdate_start.Query, node)
	if err != nil {
		return err
	}

	xp.Unitid.Query = "//archdesc[@level='collection']/did/unitid"
	xp.Unitid.Values, err = parseValues(xp.Unitid.Query, node)
	if err != nil {
		return err
	}

	xp.Unittitle.Query = "//archdesc[@level='collection']/did/unittitle"
	xp.Unittitle.Values, err = parseValues(xp.Unittitle.Query, node)
	if err != nil {
		return err
	}

	return nil
}

func parseValues(query string, node types.Node) ([]string, error) {
	return []string{}, nil
}
